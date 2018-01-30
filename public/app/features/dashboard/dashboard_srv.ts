import coreModule from 'app/core/core_module';
import { DashboardModel } from './dashboard_model';

export class DashboardSrv {
  dash: any;

  /** @ngInject */
  constructor(private backendSrv, private $rootScope, private $location) {}

  create(dashboard, meta) {
    return new DashboardModel(dashboard, meta);
  }

  setCurrent(dashboard) {
    this.dash = dashboard;
  }

  getCurrent() {
    return this.dash;
  }

  handleSaveDashboardError(clone, err) {
    if (err.data && err.data.status === 'version-mismatch') {
      err.isHandled = true;

      this.$rootScope.appEvent('confirm-modal', {
        title: 'خطا',
        text: 'شخص دیگری داشبورد را تغییر داده است',
        text2: 'همچنان می‌خواهید داشبورد را ذخیره نمایید؟',
        yesText: 'ذخیره و در نظر نگرفتن تغییرات',
        icon: 'fa-warning',
        onConfirm: () => {
          this.save(clone, { overwrite: true });
        },
      });
    }

    if (err.data && err.data.status === 'name-exists') {
      err.isHandled = true;

      this.$rootScope.appEvent('confirm-modal', {
        title: 'خطا',
        text: 'داشبوردی با همین نام قبلا ذخیره شده است.',
        text2: 'هنوز هم مایل به ذخیره کردن این داشبورد هستید؟',
        yesText: 'تغییر داشبورد قبلی',
        icon: 'fa-warning',
        onConfirm: () => {
          this.save(clone, { overwrite: true });
        },
      });
    }

    if (err.data && err.data.status === 'plugin-dashboard') {
      err.isHandled = true;

      this.$rootScope.appEvent('confirm-modal', {
        title: 'داشبورد پلاگین',
        text: err.data.message,
        text2: 'وقتی پلاگین را بروز کنید تغییرات از بین میروند. از ذخیره نسخه برای ایجاد نسخه دلخواه استفاده کنید.',
        yesText: 'Overwrite',
        icon: 'fa-warning',
        altActionText: 'ذخیره نسخه',
        onAltAction: () => {
          this.showSaveAsModal();
        },
        onConfirm: () => {
          this.save(clone, { overwrite: true });
        },
      });
    }
  }

  postSave(clone, data) {
    this.dash.version = data.version;

    var dashboardUrl = '/dashboard/db/' + data.slug;
    if (dashboardUrl !== this.$location.path()) {
      this.$location.url(dashboardUrl);
    }

    this.$rootScope.appEvent('dashboard-saved', this.dash);
    this.$rootScope.appEvent('alert-success', ['داشبورد ذخیره شد']);

    return this.dash;
  }

  save(clone, options) {
    options = options || {};
    options.folderId = this.dash.meta.folderId;

    return this.backendSrv
      .saveDashboard(clone, options)
      .then(this.postSave.bind(this, clone))
      .catch(this.handleSaveDashboardError.bind(this, clone));
  }

  saveDashboard(options, clone) {
    if (clone) {
      this.setCurrent(this.create(clone, this.dash.meta));
    }

    if (!this.dash.meta.canSave && options.makeEditable !== true) {
      return Promise.resolve();
    }

    if (this.dash.title === 'New dashboard') {
      return this.showSaveAsModal();
    }

    if (this.dash.version > 0) {
      return this.showSaveModal();
    }

    return this.save(this.dash.getSaveModelClone(), options);
  }

  showSaveAsModal() {
    this.$rootScope.appEvent('show-modal', {
      templateHtml: '<save-dashboard-as-modal dismiss="dismiss()"></save-dashboard-as-modal>',
      modalClass: 'modal--narrow',
    });
  }

  showSaveModal() {
    this.$rootScope.appEvent('show-modal', {
      templateHtml: '<save-dashboard-modal dismiss="dismiss()"></save-dashboard-modal>',
      modalClass: 'modal--narrow',
    });
  }

  starDashboard(dashboardId, isStarred) {
    let promise;

    if (isStarred) {
      promise = this.backendSrv.delete('/api/user/stars/dashboard/' + dashboardId).then(() => {
        return false;
      });
    } else {
      promise = this.backendSrv.post('/api/user/stars/dashboard/' + dashboardId).then(() => {
        return true;
      });
    }

    return promise.then(res => {
      if (this.dash && this.dash.id === dashboardId) {
        this.dash.meta.isStarred = res;
      }
      return res;
    });
  }
}

coreModule.service('dashboardSrv', DashboardSrv);
