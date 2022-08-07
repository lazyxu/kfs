import logo from './logo.svg';
import './App.css';

import { list } from './rpc/ws';

function App() {
  return (
    <div className="App">
      <header className="App-header">
        <img src={logo} className="App-logo" alt="logo" />
        <p>
          Edit <code>src/App.js</code> and save to reload.
        </p>
        <button
          className="App-link"
          href="https://reactjs.org"
          target="_blank"
          rel="noopener noreferrer"
          onClick={list}
        >
          List
        </button>
      </header>
    </div>
  );
}

export default App;
