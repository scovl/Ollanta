let _renderView = () => {};
let _showToast = () => {};

export function getRenderView() { return _renderView; }
export function getShowToast() { return _showToast; }

export function configureAdminFeature({ render, showToast }) {
  _renderView = render;
  _showToast = showToast;
}
