import { FolderPageLoader } from './folder_page_loader';
import appEvents from 'app/core/app_events';

export class FolderSettingsCtrl {
  folderPageLoader: FolderPageLoader;
  navModel: any;
  folderId: number;
  canSave = false;
  dashboard: any;
  meta: any;
  title: string;
  hasChanged: boolean;

  /** @ngInject */
  constructor(private backendSrv, navModelSrv, private $routeParams, private $location) {
    if (this.$routeParams.folderId && this.$routeParams.slug) {
      this.folderId = $routeParams.folderId;

      this.folderPageLoader = new FolderPageLoader(this.backendSrv, this.$routeParams);
      this.folderPageLoader.load(this, this.folderId, 'manage-folder-settings').then(result => {
        this.dashboard = result.dashboard;
        this.meta = result.meta;
        this.canSave = result.meta.canSave;
        this.title = this.dashboard.title;
      });
    }
  }

  save() {
    this.titleChanged();

    if (!this.hasChanged) {
      return;
    }

    this.dashboard.title = this.title.trim();

    return this.backendSrv
      .saveDashboard(this.dashboard, { overwrite: false })
      .then(result => {
        var folderUrl = this.folderPageLoader.createFolderUrl(this.folderId, this.meta.type, result.slug);
        if (folderUrl !== this.$location.path()) {
          this.$location.url(folderUrl + '/settings');
        }

        appEvents.emit('dashboard-saved');
        appEvents.emit('alert-success', ['پوشه ذخیره شد']);
      })
      .catch(this.handleSaveFolderError);
  }

  titleChanged() {
    this.hasChanged = this.dashboard.title.toLowerCase() !== this.title.trim().toLowerCase();
  }

  delete(evt) {
    if (evt) {
      evt.stopPropagation();
      evt.preventDefault();
    }

    appEvents.emit('confirm-modal', {
      title: 'حذف',
      text: `آیا می‌خواهید این پوشه و تمام داشبورد های آن را حذف کنید؟`,
      icon: 'fa-trash',
      yesText: 'حذف',
      onConfirm: () => {
        return this.backendSrv.deleteDashboard(this.meta.slug).then(() => {
          appEvents.emit('alert-success', ['Folder Deleted', `${this.dashboard.title} has been deleted`]);
          this.$location.url('dashboards');
        });
      },
    });
  }

  handleSaveFolderError(err) {
    if (err.data && err.data.status === 'version-mismatch') {
      err.isHandled = true;

      appEvents.emit('confirm-modal', {
        title: 'خطا',
        text: 'شخص دیگری پوشه را تغییر داده است.',
        text2: 'همچنان میخواهید تغییرات خود را ذخیره کنید؟',
        yesText: 'ذخیره و بازنویسی',
        icon: 'fa-warning',
        onConfirm: () => {
          this.backendSrv.saveDashboard(this.dashboard, { overwrite: true });
        },
      });
    }

    if (err.data && err.data.status === 'name-exists') {
      err.isHandled = true;

      appEvents.emit('alert-error', ['یک پوشه یا داشبورد با همین نام از قبل وجود دارد.']);
    }
  }
}
