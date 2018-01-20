import coreModule from '../../core_module';
import appEvents from 'app/core/app_events';

export class HelpCtrl {
  tabIndex: any;
  shortcuts: any;

  /** @ngInject */
  constructor() {
    this.tabIndex = 0;
    this.shortcuts = {
      Global: [
        { keys: ['g', 'h'], description: 'برو به داشبورد خانه' },
        { keys: ['g', 'p'], description: 'برو به پروفایل' },
        { keys: ['s', 'o'], description: 'بازکردن پنل جستجو' },
        { keys: ['s', 's'], description: 'Open search with starred filter' },
        { keys: ['s', 't'], description: 'Open search in tags view' },
        { keys: ['esc'], description: 'خروج از حالت تغییر' },
      ],
      Dashboard: [
        { keys: ['mod+s'], description: 'ذخیره داشبورد' },
        { keys: ['mod+h'], description: 'پنهان نمودن کنترل های سطری' },
        { keys: ['d', 'r'], description: 'Refresh تمام پنل ها' },
        { keys: ['d', 's'], description: 'تنظیمات داشبورد' },
        { keys: ['d', 'v'], description: 'Toggle in-active / view mode' },
        { keys: ['d', 'k'], description: 'Toggle kiosk mode (hides top nav)' },
        { keys: ['d', 'E'], description: 'باز کردن تمام سطرها' },
        { keys: ['d', 'C'], description: 'بستن تمام سطرها' },
        { keys: ['mod+o'], description: 'Toggle shared graph crosshair' },
      ],
      'Focused Panel': [
        { keys: ['e'], description: 'Toggle panel edit view' },
        { keys: ['v'], description: 'Toggle panel fullscreen view' },
        { keys: ['p', 's'], description: 'Open Panel Share Modal' },
        { keys: ['p', 'r'], description: 'حذف پنل' },
      ],
      'Focused Row': [{ keys: ['r', 'c'], description: 'بستن سطر' }, { keys: ['r', 'r'], description: 'حذف سطر' }],
      'Time Range': [
        { keys: ['t', 'z'], description: 'باز کردن بازه زمانی' },
        {
          keys: ['t', '<i class="fa fa-long-arrow-left"></i>'],
          description: 'حرکت دادن بازه زمانی به عقب',
        },
        {
          keys: ['t', '<i class="fa fa-long-arrow-right"></i>'],
          description: 'حرکت دادن بازه زمانی به جلو',
        },
      ],
    };
  }

  dismiss() {
    appEvents.emit('hide-modal');
  }
}

export function helpModal() {
  return {
    restrict: 'E',
    templateUrl: 'public/app/core/components/help/help.html',
    controller: HelpCtrl,
    bindToController: true,
    transclude: true,
    controllerAs: 'ctrl',
    scope: {},
  };
}

coreModule.directive('helpModal', helpModal);
