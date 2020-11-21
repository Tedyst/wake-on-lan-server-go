import React, {useState, useEffect}  from 'react';

function App() {
  const [waiting, setWaiting] = useState(0)

  const queryString = window.location.search;
  const urlParams = new URLSearchParams(queryString);
  const returnUrl = urlParams.get('returnUrl');
  
  useEffect(() => {
    if(waiting === -1)
      return;
    setTimeout(() => {
      fetch("/ping" + queryString)
      .then(res => res.json())
      .then(
        (result) => {
          if(result.success)
            setWaiting(-1);
          else setWaiting(waiting + 1);
        }
      )
    }, 1000)
  }, [waiting, queryString])

  useEffect(() => {
    fetch("/wake" + queryString)
  }, [queryString])

  if(waiting === -1)
  return (<div>Redirecting to {returnUrl}</div>)

  return (<div>Waking up server...</div>);
}

export default App;
