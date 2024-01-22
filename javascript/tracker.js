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
      visitFrom: document.referrer || 'Direct visit'
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
