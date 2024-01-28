function isNewUser() {
    const visited = localStorage.getItem('visited');
    if (visited === null) {
      localStorage.setItem('visited', 'true');
      return "new";
    } else {
      return "returning"; 
    }
  }
let visitStartTime;
function startVisitTracking() {
  visitStartTime = new Date();
}
function setVisitTimestamp() {
  const currentTimestamp = new Date().getTime();
  localStorage.setItem('lastVisitTimestamp', currentTimestamp.toString());
}
function hasVisitedWithinLast12Hours() {
  const lastVisitTimestamp = localStorage.getItem('lastVisitTimestamp');
  if (lastVisitTimestamp) {
    const currentTime = new Date().getTime();
    const lastVisitTime = parseInt(lastVisitTimestamp, 10);
    const sixHoursInMillis = 6 * 60 * 60 * 1000;

    return currentTime - lastVisitTime <= sixHoursInMillis;
  }
  return false; 
}
function getBrowserName() {
    var userAgent = navigator.userAgent;
    if (userAgent.indexOf('Chrome') !== -1) {
        return 'Google Chrome';
    } else if (userAgent.indexOf('Firefox') !== -1) {
        return 'Mozilla Firefox';
    } else if (userAgent.indexOf('Safari') !== -1) {
        return 'Apple Safari';
    } else if (userAgent.indexOf('Edge') !== -1) {
        return 'Microsoft Edge';
    } else if (userAgent.indexOf('Opera') !== -1 || userAgent.indexOf('OPR') !== -1) {
        return 'Opera';
    } else if (userAgent.indexOf('MSIE') !== -1 || userAgent.indexOf('Trident/') !== -1) {
        return 'Internet Explorer';
    } else {
        return 'Unknown Browser';
    }
}
function detectDeviceType() {
    var isMobile = /iPhone|iPad|iPod|Android|webOS|BlackBerry|IEMobile|Opera Mini/i.test(navigator.userAgent);
    var isTablet = /iPad|Android|Tablet/i.test(navigator.userAgent);

    if (isTablet) {
    return "Tablet"
    } else if (isMobile) {
        console.log('Mobile');
    return "Mobile"
    } else {
    return "Desktop"
    }
}

function detectOperatingSystem() {
    var userAgent = navigator.userAgent;
    if (/Windows/.test(userAgent)) {
    return "Windows"
    } else if (/Mac OS|Macintosh/.test(userAgent)) {
    return "Mac OS"
    } else if (/Linux/.test(userAgent)) {
    return "Linux"
    } else {
    return "Unknown"
    }
}

function analytics(domainID) {
     if (hasVisitedWithinLast12Hours()) {
      return
    } 
    let visitDuration = 0
    if (visitStartTime) {
      const visitEndTime = new Date();
      const floatValue = (visitEndTime - visitStartTime) / 1000;
      visitDuration = ~~floatValue;
    }
     const serverURL = "https://analytics-derp.koyeb.app/v1/visit" 
     // const serverURL = "http://localhost:8000/v1/visit"
    navigator.sendBeacon(serverURL,JSON.stringify({
      status:isNewUser(),
      visitDuration:visitDuration,
      domain:domainID,
      visitFrom: document.referrer || 'Direct visit',
      browser:getBrowserName(),
      device:detectDeviceType(),
      os: detectOperatingSystem()
    })
  )
  setVisitTimestamp();
}
window.addEventListener('load', startVisitTracking);
document.addEventListener('visibilitychange', function logData() {
  if (document.visibilityState === 'hidden') {
    analytics(dID)
  }
});
