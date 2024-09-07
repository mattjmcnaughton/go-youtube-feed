import "@mantine/core/styles.css";

import { MantineProvider } from "@mantine/core";

import FeedView from "./components/FeedView";

function App() {
  return (
    <MantineProvider>
      {
        <>
          <FeedView />
        </>
      }
    </MantineProvider>
  );
}

export default App;
