
import { createTheme, ThemeProvider } from '@mui/material/styles';
import App from "./App";

const theme = createTheme({
    colorSchemeSelector: 'class',
    cssVariables: true,
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

export default function () {
    return (
        <ThemeProvider theme={theme} >
            <App />
        </ThemeProvider>
    );
}
