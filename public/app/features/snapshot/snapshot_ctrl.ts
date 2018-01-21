import angular from 'angular';
import _ from 'lodash';

export class SnapshotsCtrl {
  navModel: any;
  snapshots: any;

  /** @ngInject */
  constructor(private $rootScope, private backendSrv, navModelSrv) {
    this.navModel = navModelSrv.getNav('dashboards', 'snapshots', 0);
    this.backendSrv.get('/api/dashboard/snapshots').then(result => {
      this.snapshots = result;
    });
  }

  removeSnapshotConfirmed(snapshot) {
    _.remove(this.snapshots, { key: snapshot.key });
    this.backendSrv.get('/api/snapshots-delete/' + snapshot.deleteKey).then(
      () => {
        this.$rootScope.appEvent('alert-success', ['Snapshot حذف شد', '']);
      },
      () => {
        this.$rootScope.appEvent('alert-error', ['snapshot حذف نشد', '']);
        this.snapshots.push(snapshot);
      }
    );
  }

  removeSnapshot(snapshot) {
    this.$rootScope.appEvent('confirm-modal', {
      title: 'حذف',
      text: 'آیا از حذف snapshot' + snapshot.name + 'مطمئن هستید؟',
      yesText: 'حذف',
      icon: 'fa-trash',
      onConfirm: () => {
        this.removeSnapshotConfirmed(snapshot);
      },
    });
  }
}

angular.module('grafana.controllers').controller('SnapshotsCtrl', SnapshotsCtrl);
