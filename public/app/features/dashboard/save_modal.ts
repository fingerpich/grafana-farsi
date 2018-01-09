import coreModule from 'app/core/core_module';

const template = `
<div class="modal-body">
  <div class="modal-header">
    <h2 class="modal-header-title">
      <i class="fa fa-save"></i>
      <span class="p-l-1">ذخیره کردن تغییرات</span>
    </h2>

    <a class="modal-header-close" ng-click="ctrl.dismiss();">
      <i class="fa fa-remove"></i>
    </a>
  </div>

  <form name="ctrl.saveForm" ng-submit="ctrl.save()" class="modal-content" novalidate>
    <h6 class="text-center">توضیحی در مورد تغییرات جدید بنویسید.</h6>
    <div class="p-t-2">
      <div class="gf-form">
        <label class="gf-form-hint">
          <input
            type="text"
            name="message"
            class="gf-form-input"
            placeholder="بروزرسانیه &hellip;"
            give-focus="true"
            ng-model="ctrl.message"
            ng-model-options="{allowInvalid: true}"
            ng-maxlength="this.max"
            autocomplete="off" />
          <small class="gf-form-hint-text muted" ng-cloak>
            <span ng-class="{'text-error': ctrl.saveForm.message.$invalid && ctrl.saveForm.message.$dirty }">
              {{ctrl.message.length || 0}}
            </span>
            / {{ctrl.max}} حرف
          </small>
        </label>
      </div>
    </div>

    <div class="gf-form-button-row text-center">
      <button type="submit" class="btn btn-success" ng-disabled="ctrl.saveForm.$invalid">ذخیره</button>
      <button class="btn btn-inverse" ng-click="ctrl.dismiss();">بیخیال</button>
    </div>
  </form>
</div>
`;

export class SaveDashboardModalCtrl {
  message: string;
  max: number;
  saveForm: any;
  dismiss: () => void;

  /** @ngInject */
  constructor(private dashboardSrv) {
    this.message = '';
    this.max = 64;
  }

  save() {
    if (!this.saveForm.$valid) {
      return;
    }

    var dashboard = this.dashboardSrv.getCurrent();
    var saveModel = dashboard.getSaveModelClone();
    var options = { message: this.message };

    return this.dashboardSrv.save(saveModel, options).then(this.dismiss);
  }
}

export function saveDashboardModalDirective() {
  return {
    restrict: 'E',
    template: template,
    controller: SaveDashboardModalCtrl,
    bindToController: true,
    controllerAs: 'ctrl',
    scope: { dismiss: '&' },
  };
}

coreModule.directive('saveDashboardModal', saveDashboardModalDirective);
