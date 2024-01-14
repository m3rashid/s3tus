import '@uppy/core/dist/style.min.css';
import '@uppy/dashboard/dist/style.min.css';

import React, { useEffect } from 'react';
import Uppy from '@uppy/core';
import Dashboard from '@uppy/dashboard';

import Tus from '@uppy/tus';

function App() {
  const appRef = React.useRef(null);
  const [uppy, setUppy] = React.useState<Uppy>();

  useEffect(() => {
    if (uppy) return;
    if (appRef.current && !uppy) {
      const newUppy = new Uppy({ debug: true })
        .use(Dashboard, { inline: true, target: '#uppy-dashboard' })
        .use(Tus, { endpoint: 'http://localhost:5000/files/' });
      setUppy(newUppy);
    }
  }, []);

  return (
    <>
      <img
        height={500}
        src='http://localhost:5000/file/d346bbee91bc6dd98d93e2563113275c'
        // src='https://s3.ap-south-1.amazonaws.com/go.awesome.stack/d346bbee91bc6dd98d93e2563113275c'
      />
      {/* <video
        controls
        height={500}
        src='http://localhost:5000/file/40b011750f6893fb772ff68844997e2b'
      /> */}
      <div id='uppy-dashboard' ref={appRef} />
    </>
  );
}

export default App;
