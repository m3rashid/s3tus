import '@uppy/core/dist/style.min.css';
import '@uppy/dashboard/dist/style.min.css';

import Tus from '@uppy/tus';
import Uppy from '@uppy/core';
import { onMount } from 'solid-js';
import Dashboard from '@uppy/dashboard';

function App() {
  onMount(() => {
    new Uppy({ debug: true })
      .use(Dashboard, { inline: true, target: '#uppy-dashboard' })
      .use(Tus, { endpoint: 'http://localhost:5000/files/' });
  });

  return (
    <>
      {[
        'd346bbee91bc6dd98d93e2563113275c',
        '199764f14e0b6fbedb53452c6d37eb06',
      ].map((imageId) => (
        <img height={500} src={`http://localhost:5000/file/${imageId}`} />
      ))}
      <div id='uppy-dashboard' />
    </>
  );
}

export default App;
