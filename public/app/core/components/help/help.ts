import coreModule from '../../core_module';
import appEvents from 'app/core/app_events';

export class HelpCtrl {
  tabIndex: any;
  shortcuts: any;

  /** @ngInject */
  constructor() {
    this.tabIndex = 0;
    this.shortcuts = {
      عمومی: [
        { keys: ['g', 'h'], description: 'برو به داشبورد خانه' },
        { keys: ['g', 'p'], description: 'برو به پروفایل' },
        { keys: ['s', 'o'], description: 'بازکردن پنل جستجو' },
        { keys: ['s', 's'], description: 'باز کردن جستجو با فیلتر ستاره شده ها' },
        { keys: ['s', 't'], description: 'باز کردن جستجو در محیط تگ ها' },
        { keys: ['esc'], description: 'خروج از حالت تغییر' },
      ],
      داشبورد: [
        { keys: ['mod+s'], description: 'ذخیره داشبورد' },
        { keys: ['mod+h'], description: 'پنهان نمودن کنترل های سطری' },
        { keys: ['d', 'r'], description: 'Refresh تمام پنل ها' },
        { keys: ['d', 's'], description: 'تنظیمات داشبورد' },
        { keys: ['d', 'v'], description: 'تغییر حالت بین in-active و view mode' },
        { keys: ['d', 'k'], description: 'تغییر حالت kiosk mode بین  hides و top  و nav' },
        { keys: ['d', 'E'], description: 'باز کردن تمام سطرها' },
        { keys: ['d', 'C'], description: 'بستن تمام سطرها' },
        { keys: ['mod+o'], description: 'تغییر حالت shared graph crosshair' },
      ],
      'در پنل باز شده': [
        { keys: ['e'], description: 'تغییر حالت پنل بین edit و view' },
        { keys: ['v'], description: 'تغییر حالت پنل بین fullscreen  و view' },
        { keys: ['p', 's'], description: 'باز کردن پنل اشتراک گذاری' },
        { keys: ['p', 'r'], description: 'حذف پنل' },
      ],
      'در سطر فعال': [{ keys: ['r', 'c'], description: 'بستن سطر' }, { keys: ['r', 'r'], description: 'حذف سطر' }],
      'در بازه زمانی': [
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
