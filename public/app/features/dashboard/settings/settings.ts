import { coreModule, appEvents, contextSrv } from 'app/core/core';
import { DashboardModel } from '../dashboard_model';
import $ from 'jquery';
import _ from 'lodash';

export class SettingsCtrl {
  dashboard: DashboardModel;
  isOpen: boolean;
  viewId: string;
  json: string;
  alertCount: number;
  canSaveAs: boolean;
  canDelete: boolean;
  sections: any[];

  /** @ngInject */
  constructor(private $scope, private $location, private $rootScope, private backendSrv, private dashboardSrv) {
    // temp hack for annotations and variables editors
    // that rely on inherited scope
    $scope.dashboard = this.dashboard;

    this.$scope.$on('$destroy', () => {
      this.dashboard.updateSubmenuVisibility();
      this.$rootScope.$broadcast('refresh');
    });

    this.canSaveAs = contextSrv.isEditor;
    this.canDelete = this.dashboard.meta.canSave;

    this.buildSectionList();
    this.onRouteUpdated();

    $rootScope.onAppEvent('$routeUpdate', this.onRouteUpdated.bind(this), $scope);
  }

  buildSectionList() {
    this.sections = [];

    if (this.dashboard.meta.canEdit) {
      this.sections.push({
        title: 'عمومی',
        id: 'settings',
        icon: 'gicon gicon-preferences',
      });
      this.sections.push({
        title: 'حاشیه‌ها',
        id: 'annotations',
        icon: 'gicon gicon-annotation',
      });
      this.sections.push({
        title: 'متغیرها',
        id: 'templating',
        icon: 'gicon gicon-variable',
      });
      this.sections.push({
        title: 'لینک‌ ها',
        id: 'links',
        icon: 'gicon gicon-link',
      });
    }

    if (this.dashboard.id && this.dashboard.meta.canSave) {
      this.sections.push({
        title: 'نسخه ها',
        id: 'versions',
        icon: 'fa fa-fw fa-history',
      });
    }

    if (this.dashboard.meta.canMakeEditable) {
      this.sections.push({
        title: 'عمومی',
        icon: 'gicon gicon-preferences',
        id: 'make_editable',
      });
    }

    this.sections.push({
      title: 'نمایش JSON',
      id: 'view_json',
      icon: 'gicon gicon-json',
    });

    const params = this.$location.search();
    const url = this.$location.path();

    for (let section of this.sections) {
      const sectionParams = _.defaults({ editview: section.id }, params);
      section.url = url + '?' + $.param(sectionParams);
    }
  }

  onRouteUpdated() {
    this.viewId = this.$location.search().editview;

    if (this.viewId) {
      this.json = JSON.stringify(this.dashboard.getSaveModelClone(), null, 2);
    }

    if (this.viewId === 'settings' && this.dashboard.meta.canMakeEditable) {
      this.viewId = 'make_editable';
    }

    const currentSection = _.find(this.sections, { id: this.viewId });
    if (!currentSection) {
      this.sections.unshift({
        title: 'یافت نشد',
        id: '404',
        icon: 'fa fa-fw fa-warning',
      });
      this.viewId = '404';
    }
  }

  openSaveAsModal() {
    this.dashboardSrv.showSaveAsModal();
  }

  hideSettings() {
    var urlParams = this.$location.search();
    delete urlParams.editview;
    setTimeout(() => {
      this.$rootScope.$apply(() => {
        this.$location.search(urlParams);
      });
    });
  }

  makeEditable() {
    this.dashboard.editable = true;
    this.dashboard.meta.canMakeEditable = false;
    this.dashboard.meta.canEdit = true;
    this.dashboard.meta.canSave = true;
    this.canDelete = true;
    this.viewId = 'settings';
    this.buildSectionList();

    const currentSection = _.find(this.sections, { id: this.viewId });
    this.$location.url(currentSection.url);
  }

  deleteDashboard() {
    var confirmText = '';
    var text2 = this.dashboard.title;

    const alerts = _.sumBy(this.dashboard.panels, panel => {
      return panel.alert ? 1 : 0;
    });

    if (alerts > 0) {
      confirmText = 'حذف';
      text2 = `این داشبورد شامل  ${alerts} هشدار است. با حذف این داشبورد شما تمام این هشدارها را نیز حذف خواهید نمود.`;
    }

    appEvents.emit('confirm-modal', {
      title: 'حذف',
      text: 'آیا می‌خواهید این داشبورد را حذف نمایید؟',
      text2: text2,
      icon: 'fa-trash',
      confirmText: confirmText,
      yesText: 'حذف',
      onConfirm: () => {
        this.dashboard.meta.canSave = false;
        this.deleteDashboardConfirmed();
      },
    });
  }

  deleteDashboardConfirmed() {
    this.backendSrv.deleteDashboard(this.dashboard.meta.slug).then(() => {
      appEvents.emit('alert-success', ['داشبورد حذف شد', this.dashboard.title + ' حذف شده است.']);
      this.$location.url('/');
    });
  }

  onFolderChange(folder) {
    this.dashboard.meta.folderId = folder.id;
    this.dashboard.meta.folderTitle = folder.title;
  }
}

export function dashboardSettings() {
  return {
    restrict: 'E',
    templateUrl: 'public/app/features/dashboard/settings/settings.html',
    controller: SettingsCtrl,
    bindToController: true,
    controllerAs: 'ctrl',
    transclude: true,
    scope: { dashboard: '=' },
  };
}

coreModule.directive('dashboardSettings', dashboardSettings);
