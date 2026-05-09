import { state } from '../core/state.js';

export { configureAdminFeature } from './admin/context.js';

export { loadGateData, renderGateTab, bindGatesContent } from './admin/gates.js';
export { loadWebhooksData, renderWebhooksTab, bindWebhooksContent } from './admin/webhooks.js';
export { loadProfilesData, renderProfilesTab, bindProfilesContent } from './admin/profiles.js';
export { loadCustomRulesData, renderCustomRulesTab, bindCustomRulesContent } from './admin/custom-rules.js';
export { loadAIProvidersData, renderAIProvidersTab, bindAIProvidersContent } from './admin/ai-providers.js';
export {
  loadBackgroundTasksData,
  loadBackgroundTaskDetail,
  renderAdminLinksTab,
  renderBackgroundTasksPage,
  bindBackgroundTasksContent,
} from './admin/background-tasks.js';

import { bindBackgroundTasksContent } from './admin/background-tasks.js';
import { bindGatesContent } from './admin/gates.js';
import { bindWebhooksContent } from './admin/webhooks.js';
import { bindProfilesContent } from './admin/profiles.js';
import { bindCustomRulesContent } from './admin/custom-rules.js';
import { bindAIProvidersContent } from './admin/ai-providers.js';

export function bindAdminTabContent() {
  bindBackgroundTasksContent();
  if (!state.currentProject) return;
  bindGatesContent();
  bindWebhooksContent();
  bindProfilesContent();
  bindCustomRulesContent();
  bindAIProvidersContent();
}
