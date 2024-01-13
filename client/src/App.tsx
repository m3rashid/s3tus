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

  return <div id='uppy-dashboard' ref={appRef} />;
}

export default App;
