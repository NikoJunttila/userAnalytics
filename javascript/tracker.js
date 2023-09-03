// Function to check if the user is a returning or new user based on local storage
  function isNewUser() {
    const visited = localStorage.getItem('visited');
    if (visited === null) {
      localStorage.setItem('visited', 'true');
      return "new"; // New user
    } else {
      return "returning"; // Returning user
    }
  }
  // Function to get the referring URL (where the user came from)
  function getReferringURL() {
    return document.referrer || 'Direct visit';
  }
  let visitStartTime;
function startVisitTracking() {
  visitStartTime = new Date();
}
function toInt32(number) {
  return number | 0;
}
function logVisitDuration() {
  if (visitStartTime) {
    const visitEndTime = new Date();
    const visitDuration = visitEndTime - visitStartTime;

    console.log(`User visited for ${visitDuration / 1000} seconds`);
    visitStartTime = null;
  }
}
  // Main function to collect and send user data
  function collectAndSendUserData() {
    //const country = await getUserCountry();
    const isNew = isNewUser();
    const referringURL = getReferringURL();
    let visitDuration = 0
    if (visitStartTime) {
      const visitEndTime = new Date();
      floatValue = (visitEndTime - visitStartTime) / 1000;
      const visitDuration = toInt32(floatValue);
  
      console.log(`User visited for ${visitDuration} seconds`);
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
        domain:"22b567c5-2736-436d-81f7-edab07857f02",
        visitFrom:referringURL
      })
  })
}

window.addEventListener('load', startVisitTracking);
window.addEventListener('beforeunload', collectAndSendUserData);

/* collectAndSendUserData(); */
  