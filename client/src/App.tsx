import '@uppy/core/dist/style.min.css';
import '@uppy/dashboard/dist/style.min.css';

import Tus from '@uppy/tus';
import Uppy from '@uppy/core';
import { createSignal, onMount } from 'solid-js';
import Dashboard from '@uppy/dashboard';

const App = () => {
  const [getImages, setImages] = createSignal<string[]>([]);

  const handleGetImages = async () => {
    const res = await fetch('http://localhost:5000/auth/', {
      method: 'GET',
      credentials: 'include',
    });
    if (!res.ok) throw new Error('Bad response');
    setImages([
      'd346bbee91bc6dd98d93e2563113275c',
      '199764f14e0b6fbedb53452c6d37eb06',
    ]);
  };

  onMount(() => {
    handleGetImages().catch(console.log);
    new Uppy({ debug: true })
      .use(Dashboard, { inline: true, target: '#uppy-dashboard' })
      .use(Tus, { endpoint: 'http://localhost:5000/files/' });
  });

  return (
    <>
      {getImages().map((imageId) => (
        <img height={500} src={`http://localhost:5000/file/${imageId}`} />
      ))}
      <div id='uppy-dashboard' />
    </>
  );
};

export default App;
