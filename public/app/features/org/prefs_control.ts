import config from 'app/core/config';
import coreModule from 'app/core/core_module';

export class PrefsControlCtrl {
  prefs: any;
  oldTheme: any;
  prefsForm: any;
  mode: string;

  timezones: any = [
    { value: '', text: 'پیش فرض' },
    { value: 'browser', text: 'زمان بروزر شما' },
    { value: 'utc', text: 'UTC' },
  ];
  themes: any = [{ value: '', text: 'پیش فرض' }, { value: 'dark', text: 'تاریک' }, { value: 'light', text: 'روشن' }];

  /** @ngInject **/
  constructor(private backendSrv, private $location) {}

  $onInit() {
    return this.backendSrv.get(`/api/${this.mode}/preferences`).then(prefs => {
      this.prefs = prefs;
      this.oldTheme = prefs.theme;
    });
  }

  updatePrefs() {
    if (!this.prefsForm.$valid) {
      return;
    }

    var cmd = {
      theme: this.prefs.theme,
      timezone: this.prefs.timezone,
      homeDashboardId: this.prefs.homeDashboardId,
    };

    this.backendSrv.put(`/api/${this.mode}/preferences`, cmd).then(() => {
      window.location.href = config.appSubUrl + this.$location.path();
    });
  }
}

var template = `
<form name="ctrl.prefsForm" class="section gf-form-group">
  <h3 class="page-heading">اولویت ها</h3>

  <div class="gf-form">
    <span class="gf-form-label width-11">UI تم</span>
    <div class="gf-form-select-wrapper max-width-20">
      <select class="gf-form-input" ng-model="ctrl.prefs.theme" ng-options="f.value as f.text for f in ctrl.themes"></select>
    </div>
  </div>

  <div class="gf-form">
    <span class="gf-form-label width-11">
      داشبورد صفحه اول
      <info-popover mode="right-normal">
        داشبورد مورد نظر پیدا نشد؟ ابتدا داشبورد مورد نظر را ستاره کنید و سپس در این بخش ظاهر میگردد
      </info-popover>
    </span>
    <dashboard-selector class="gf-form-select-wrapper max-width-20" model="ctrl.prefs.homeDashboardId">
    </dashboard-selector>
  </div>

  <div class="gf-form">
    <label class="gf-form-label width-11">موقعیت زمانی</label>
    <div class="gf-form-select-wrapper max-width-20">
      <select class="gf-form-input" ng-model="ctrl.prefs.timezone" ng-options="f.value as f.text for f in ctrl.timezones"></select>
    </div>
  </div>

  <div class="gf-form-button-row">
    <button type="submit" class="btn btn-success" ng-click="ctrl.updatePrefs()">ذخیره</button>
  </div>
</form>
`;

export function prefsControlDirective() {
  return {
    restrict: 'E',
    controller: PrefsControlCtrl,
    bindToController: true,
    controllerAs: 'ctrl',
    template: template,
    scope: {
      mode: '@',
    },
  };
}

coreModule.directive('prefsControl', prefsControlDirective);
