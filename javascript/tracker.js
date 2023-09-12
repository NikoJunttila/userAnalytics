function isNewUser() {
  const visited = localStorage.getItem('visited');
  if (visited === null) {
    localStorage.setItem('visited', 'true');
    return "new";
  } else {
    return "returning"; 
  }
}
function getReferringURL() {
  return document.referrer || 'direct';
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
  const twelveHoursInMillis = 12 * 60 * 60 * 1000;

  return currentTime - lastVisitTime <= twelveHoursInMillis;
}
return false; 
}
function analytics(domainID) {
  if (hasVisitedWithinLast12Hours()) {
    return
  }
  const isNew = isNewUser();
  const referringURL = getReferringURL();
  let visitDuration = 0
  if (visitStartTime) {
    const visitEndTime = new Date();
    const floatValue = (visitEndTime - visitStartTime) / 1000;
    visitDuration = ~~floatValue;
    visitStartTime = null;
  }
  // Send the data to the server using an HTTP POST request
  fetch('http://localhost:8000/v1/visit', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      status:isNew,
      visitDuration:visitDuration ,
      domain:domainID,
      visitFrom:referringURL
    })
})
setVisitTimestamp();
}
window.addEventListener('load', startVisitTracking);
window.onpagehide = (event) => {
if (event.persisted) {
  localStorage.removeItem('lastVisitTimestamp')
} else{
  analytics(dID)
}
};
