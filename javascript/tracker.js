function isNewUser() {
    const visited = localStorage.getItem("vis123");
    if (visited === null) {
      localStorage.setItem("vis123", "true");
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
    localStorage.setItem("lvt123", currentTimestamp.toString());
  }
  function hasVisitedWithinLast12Hours() {
    const lastVisitTimestamp = localStorage.getItem("lvt123");
    if (lastVisitTimestamp) {
      const currentTime = new Date().getTime();
      const lastVisitTime = parseInt(lastVisitTimestamp, 10);
      const sixHoursInMillis = 6 * 60 * 60 * 1000;
  
      return currentTime - lastVisitTime <= sixHoursInMillis;
    }
    return false;
  }
  function testLocal() {
    const currentUrl = window.location.href;
    if (currentUrl.includes("localhost")) {
      return true;
    } else {
      return false;
    }
  }
  function getBrowserName() {
    if (
      (navigator.userAgent.indexOf("Opera") ||
        navigator.userAgent.indexOf("OPR")) != -1
    ) {
      return "Opera";
    } else if (navigator.userAgent.indexOf("Edg") != -1) {
      return "Edge";
    } else if (navigator.userAgent.indexOf("Chrome") != -1) {
      return "Chrome";
    } else if (navigator.userAgent.indexOf("Safari") != -1) {
      return "Safari";
    } else if (navigator.userAgent.indexOf("Firefox") != -1) {
      return "Firefox";
    } else if (
      navigator.userAgent.indexOf("MSIE") != -1 ||
      !!document.documentMode == true
    ) {
      //IF IE > 10
      return "IE";
    } else {
      return "unknown";
    }
  }
  function detectDeviceType() {
    var isMobile =
      /iPhone|iPad|iPod|Android|webOS|BlackBerry|IEMobile|Opera Mini/i.test(
        navigator.userAgent
      );
    var isTablet = /iPad|Android|Tablet/i.test(navigator.userAgent);
  
    if (isMobile) {
      return "Mobile";
    } else if (isTablet) {
      return "Tablet";
    } else {
      return "Desktop/Laptop";
    }
  }
  
  function getOS() {
    const userAgent = window.navigator.userAgent,
      platform =
        window.navigator?.userAgentData?.platform || window.navigator.platform,
      macosPlatforms = ["macOS", "Macintosh", "MacIntel", "MacPPC", "Mac68K"],
      windowsPlatforms = ["Win32", "Win64", "Windows", "WinCE"],
      iosPlatforms = ["iPhone", "iPad", "iPod"];
    let os = "unknown";
  
    if (macosPlatforms.indexOf(platform) !== -1) {
      os = "Mac";
    } else if (iosPlatforms.indexOf(platform) !== -1) {
      os = "iOS";
    } else if (windowsPlatforms.indexOf(platform) !== -1) {
      os = "Windows";
    } else if (/Android/.test(userAgent)) {
      os = "Android";
    } else if (/Linux/.test(platform)) {
      os = "Linux";
    }
    return os;
  }
  const serverURL = "https://analytics-derp.koyeb.app/v1/";
  function analytics(domainID) {
    if (testLocal() || hasVisitedWithinLast12Hours()) {
      return;
    }
    let visitDuration = 0;
    if (visitStartTime) {
      const visitEndTime = new Date();
      const floatValue = (visitEndTime - visitStartTime) / 1000;
      visitDuration = ~~floatValue;
    }
    navigator.sendBeacon(
      `${serverURL}visit`,
      JSON.stringify({
        status: isNewUser(),
        visitDuration: visitDuration,
        domain: domainID,
        visitFrom: document.referrer || "Direct visit",
        browser: getBrowserName(),
        device: detectDeviceType(),
        os: getOS(),
      })
    );
    setVisitTimestamp();
  }
  window.addEventListener("load", startVisitTracking);
  document.addEventListener("visibilitychange", function logData() {
    if (document.visibilityState === "hidden") {
      analytics(dID);
    }
  });

  let previousUrl = '';

  const observer = new MutationObserver(function(mutations) {
    if (location.href !== previousUrl) {
    if (testLocal())return
      previousUrl = location.href;
      navigator.sendBeacon(
        `${serverURL}pageVisit`,
        JSON.stringify({
          domain: dID,
          page: location.href
        })
      );
    }
  });
  
  const config = { subtree: true, childList: true };
  observer.observe(document.body, config);
