(() => {
  document.documentElement.addEventListener('touchstart', (event) => {
    if (event.touches.length > 1) { event.preventDefault(); }
  }, { passive: false });

  let lastTouchEnd = 0;
  document.documentElement.addEventListener('touchend', (event) => {
    let now = Date.now();
    if (now - lastTouchEnd <= 500) { event.preventDefault() }
    lastTouchEnd = now;
  }, { passive: false });

})()

const IS_MOBILE = /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(navigator.userAgent);

const DIR_ID_PREFIX = 'game__operate-direction--';
const MOBILE_DIR_MAP = [
  { id: `${DIR_ID_PREFIX}up`, keyCode: 38 },
  { id: `${DIR_ID_PREFIX}down`, keyCode: 40 },
  { id: `${DIR_ID_PREFIX}left`, keyCode: 37 },
  { id: `${DIR_ID_PREFIX}right`, keyCode: 39 },
];