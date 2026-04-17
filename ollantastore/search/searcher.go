package search

import (
	"context"
	"fmt"
	"strings"

	meilisearch "github.com/meilisearch/meilisearch-go"
)

// SearchRequest describes a full-text search query with optional filters and facets.
type SearchRequest struct {
	Query  string
	Filter map[string]string // field → value; joined as AND filters
	Facets []string
	Sort   []string
	Limit  int
	Offset int
}

// SearchResult is the unified response from a Meilisearch search call.
type SearchResult struct {
	Hits              []map[string]interface{}   `json:"hits"`
	TotalHits         int64                      `json:"total_hits"`
	FacetDistribution map[string]map[string]int  `json:"facet_distribution"`
	ProcessingTimeMs  int64                      `json:"processing_time_ms"`
}

// MeilisearchSearcher executes full-text queries against Meilisearch indexes.
type MeilisearchSearcher struct {
	client meilisearch.ServiceManager
}

// NewMeilisearchSearcher creates a searcher connected to the given Meilisearch instance.
func NewMeilisearchSearcher(cfg IndexerConfig) (*MeilisearchSearcher, error) {
	client := meilisearch.New(cfg.Host, meilisearch.WithAPIKey(cfg.APIKey))
	return &MeilisearchSearcher{client: client}, nil
}

// SearchIssues executes a query against the issues index.
func (s *MeilisearchSearcher) SearchIssues(_ context.Context, req SearchRequest) (*SearchResult, error) {
	return s.search(indexIssues, req)
}

// SearchProjects executes a query against the projects index.
func (s *MeilisearchSearcher) SearchProjects(_ context.Context, req SearchRequest) (*SearchResult, error) {
	return s.search(indexProjects, req)
}

func (s *MeilisearchSearcher) search(index string, req SearchRequest) (*SearchResult, error) {
	limit := int64(req.Limit)
	if limit <= 0 {
		limit = 20
	}

	sr := &meilisearch.SearchRequest{
		Limit:  limit,
		Offset: int64(req.Offset),
	}

	if len(req.Filter) > 0 {
		parts := make([]string, 0, len(req.Filter))
		for k, v := range req.Filter {
			parts = append(parts, fmt.Sprintf("%s = %q", k, v))
		}
		sr.Filter = strings.Join(parts, " AND ")
	}

	if len(req.Facets) > 0 {
		sr.Facets = req.Facets
	}

	if len(req.Sort) > 0 {
		sr.Sort = req.Sort
	}

	resp, err := s.client.Index(index).Search(req.Query, sr)
	if err != nil {
		return nil, fmt.Errorf("meilisearch search %s: %w", index, err)
	}

	result := &SearchResult{
		TotalHits:        resp.EstimatedTotalHits,
		ProcessingTimeMs: resp.ProcessingTimeMs,
		Hits:             make([]map[string]interface{}, len(resp.Hits)),
	}

	for i, h := range resp.Hits {
		if m, ok := h.(map[string]interface{}); ok {
			result.Hits[i] = m
		} else {
			result.Hits[i] = map[string]interface{}{"_raw": h}
		}
	}

	if fd, ok := resp.FacetDistribution.(map[string]interface{}); ok {
		result.FacetDistribution = make(map[string]map[string]int, len(fd))
		for facetName, facetVals := range fd {
			if vals, ok := facetVals.(map[string]interface{}); ok {
				m := make(map[string]int, len(vals))
				for k, v := range vals {
					if n, ok := v.(float64); ok {
						m[k] = int(n)
					}
				}
				result.FacetDistribution[facetName] = m
			}
		}
	}

	return result, nil
}
