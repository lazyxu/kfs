import {experimental_extendTheme} from "@mui/material";

import {Experimental_CssVarsProvider as CssVarsProvider} from '@mui/material/styles';
import React, {useState} from "react";
import App from "./App";

const theme = experimental_extendTheme({
    colorSchemes: {
        light: {
            palette: {
                folder: '#78d0f9',
                file: "#39370d",
                fileViewHeader : '#d0d7de',
            },
        },
        dark: {
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
