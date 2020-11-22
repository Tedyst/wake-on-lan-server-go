import React, { useState, useEffect } from "react";

function App() {
  const [waiting, setWaiting] = useState(0);

  const queryString = window.location.search;
  const urlParams = new URLSearchParams(queryString);
  const redirectURL = urlParams.get("redirectURL");

  useEffect(() => {
    if (waiting === -1) return;
    setTimeout(() => {
      fetch("/ping" + queryString)
        .then((res) => res.json())
        .then((result) => {
          if (result.success) setWaiting(-1);
          else setWaiting(waiting + 1);
        });
    }, 1000);
  }, [waiting, queryString]);

  useEffect(() => {
    fetch("/wake" + queryString);
  }, [queryString]);

  if (waiting === -1) {
    if (redirectURL != null) window.location.replace(redirectURL);
    return (
      <div>
        Redirecting to {redirectURL}. If it does not work, the redirect has been
        setup improperly.
      </div>
    );
  }

  return <div>Waking up server...</div>;
}

export default App;
