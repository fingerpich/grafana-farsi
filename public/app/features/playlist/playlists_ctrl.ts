import _ from 'lodash';
import coreModule from '../../core/core_module';

export class PlaylistsCtrl {
  playlists: any;
  navModel: any;

  /** @ngInject */
  constructor(private $scope, private backendSrv, navModelSrv) {
    this.navModel = navModelSrv.getNav('dashboards', 'playlists', 0);

    backendSrv.get('/api/playlists').then(result => {
      this.playlists = result;
    });
  }

  removePlaylistConfirmed(playlist) {
    _.remove(this.playlists, { id: playlist.id });

    this.backendSrv.delete('/api/playlists/' + playlist.id).then(
      () => {
        this.$scope.appEvent('alert-success', ['Playlist حذف شد', '']);
      },
      () => {
        this.$scope.appEvent('alert-error', ['playlist حذف نشد', '']);
        this.playlists.push(playlist);
      }
    );
  }

  removePlaylist(playlist) {
    this.$scope.appEvent('confirm-modal', {
      title: 'حذف',
      text: 'آیا از حذف ' + playlist.name + ' مطمئن هستید؟',
      yesText: 'حذف',
      icon: 'fa-trash',
      onConfirm: () => {
        this.removePlaylistConfirmed(playlist);
      },
    });
  }
}

coreModule.controller('PlaylistsCtrl', PlaylistsCtrl);
