import angular from 'angular';

export class AdminListOrgsCtrl {
  /** @ngInject */
  constructor($scope, backendSrv, navModelSrv) {
    $scope.init = function() {
      $scope.navModel = navModelSrv.getNav('cfg', 'admin', 'global-orgs', 1);
      $scope.getOrgs();
    };

    $scope.getOrgs = function() {
      backendSrv.get('/api/orgs').then(function(orgs) {
        $scope.orgs = orgs;
      });
    };

    $scope.deleteOrg = function(org) {
      $scope.appEvent('confirm-modal', {
        title: 'حذف',
        text: 'آیا میخواهید سازمان ' + org.name + ' را حذف نمایید؟',
        text2: 'تمام داشبورد های مربوط به این سازمان حذف خواهند شد!',
        icon: 'fa-trash',
        yesText: 'حذف',
        onConfirm: function() {
          backendSrv.delete('/api/orgs/' + org.id).then(function() {
            $scope.getOrgs();
          });
        },
      });
    };

    $scope.init();
  }
}

angular.module('grafana.controllers').controller('AdminListOrgsCtrl', AdminListOrgsCtrl);
