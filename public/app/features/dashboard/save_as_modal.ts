import coreModule from 'app/core/core_module';

const template = `
<div class="modal-body">
	<div class="modal-header">
		<h2 class="modal-header-title">
			<i class="fa fa-copy"></i>
			<span class="p-l-1">ذخیره به صورت...</span>
		</h2>

		<a class="modal-header-close" ng-click="ctrl.dismiss();">
			<i class="fa fa-remove"></i>
		</a>
	</div>

	<form name="ctrl.saveForm" ng-submit="ctrl.save()" class="modal-content" novalidate>
		<div class="p-t-2">
			<div class="gf-form">
				<label class="gf-form-label width-7">نام جدید</label>
				<input type="text" class="gf-form-input" ng-model="ctrl.clone.title" give-focus="true" required>
			</div>
      <div class="gf-form">
        <folder-picker initial-folder-id="ctrl.folderId"
                       on-change="ctrl.onFolderChange($folder)"
                       enter-folder-creation="ctrl.onEnterFolderCreation()"
                       exit-folder-creation="ctrl.onExitFolderCreation()"
                       enable-create-new="true"
                       label-class="width-7">
        </folder-picker>
      </div>
		</div>

		<div class="gf-form-button-row text-center">
			<button type="submit" class="btn btn-success" ng-disabled="ctrl.saveForm.$invalid || !ctrl.isValidFolderSelection">ذخیره</button>
			<a class="btn-text" ng-click="ctrl.dismiss();">بی خیال</a>
		</div>
	</form>
</div>
`;

export class SaveDashboardAsModalCtrl {
  clone: any;
  folderId: any;
  isValidFolderSelection = true;
  dismiss: () => void;

  /** @ngInject */
  constructor(private dashboardSrv) {
    var dashboard = this.dashboardSrv.getCurrent();
    this.clone = dashboard.getSaveModelClone();
    this.clone.id = null;
    this.clone.title += ' کپی';
    this.clone.editable = true;
    this.clone.hideControls = false;
    this.folderId = dashboard.meta.folderId;

    // remove alerts if source dashboard is already persisted
    // do not want to create alert dupes
    if (dashboard.id > 0) {
      this.clone.panels.forEach(panel => {
        if (panel.type === 'graph' && panel.alert) {
          delete panel.thresholds;
        }
        delete panel.alert;
      });
    }

    delete this.clone.autoUpdate;
  }

  save() {
    return this.dashboardSrv.save(this.clone).then(this.dismiss);
  }

  onEnterFolderCreation() {
    this.isValidFolderSelection = false;
  }

  onExitFolderCreation() {
    this.isValidFolderSelection = true;
  }

  keyDown(evt) {
    if (this.isValidFolderSelection && evt.keyCode === 13) {
      this.save();
    }
  }

  onFolderChange(folder) {
    this.clone.folderId = folder.id;
  }
}

export function saveDashboardAsDirective() {
  return {
    restrict: 'E',
    template: template,
    controller: SaveDashboardAsModalCtrl,
    bindToController: true,
    controllerAs: 'ctrl',
    scope: { dismiss: '&' },
  };
}

coreModule.directive('saveDashboardAsModal', saveDashboardAsDirective);
