import { HoxRoot } from 'hox';
import ReactDOM from 'react-dom/client';

import './index.scss';
import './normalize.css';
import './reset.scss';

import reportWebVitals from './reportWebVitals';
import ThemeApp from "./ThemeApp";

window.kfs = {
    env: import.meta.env,
};

const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(
    // <React.StrictMode>
    <HoxRoot>
        <ThemeApp/>
    </HoxRoot>
    // </React.StrictMode>,
);

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();
