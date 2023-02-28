import {experimental_extendTheme} from "@mui/material";

import {Experimental_CssVarsProvider as CssVarsProvider} from '@mui/material/styles';
import React, {useState} from "react";
import App from "./App";

const theme = experimental_extendTheme({
    colorSchemes: {
        light: {
            background: {
                primary: 'rgb(255, 255, 255)',
                secondary: 'rgb(245, 245, 246)',
            },
            context: {
                primary: 'rgb(37, 38, 43)',
                secondary: 'rgba(37, 38, 43, 0.72)',
                tertiary: 'rgba(37, 38, 43, 0.36)',
                quaternary: 'rgba(37, 38, 43, 0.18)',
            },
            palette: {
                folder: '#78d0f9',
                file: "#39370d",
                fileViewHeader: '#d0d7de',
            },
        },
        dark: {
            background: {
                primary: 'rgb(17, 17, 19)',
                secondary: 'rgb(34, 34, 38)',
            },
            context: {
                primary: 'rgb(255, 255, 255)',
                secondary: 'rgba(255, 255, 255, 0.72)',
                tertiary: 'rgba(255, 255, 255, 0.36)',
                quaternary: 'rgba(255, 255, 255, 0.18)',
            },
            palette: {
                folder: '#65d1fd',
                file: "#eeeeee",
                fileViewHeader: '#333942',
            },
        },
    },
});

const useEnhancedEffect =
    typeof window !== 'undefined' ? React.useLayoutEffect : React.useEffect;

export default function () {
    // the `node` is used for attaching CSS variables to this demo, you might not need it in your application.
    const [node, setNode] = useState(null);
    useEnhancedEffect(() => {
        setNode(document.getElementById('css-vars-custom-theme'));
    }, []);
    return (
        <div id="css-vars-custom-theme">
            <CssVarsProvider
                theme={theme}
                colorSchemeNode={node || null}
                colorSchemeSelector="#css-vars-custom-theme"
                colorSchemeStorageKey="custom-theme-color-scheme"
                modeStorageKey="custom-theme-mode"
            >
                <App/>
            </CssVarsProvider>
        </div>
    );
}
