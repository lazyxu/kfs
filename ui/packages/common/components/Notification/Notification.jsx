import { Close } from "@mui/icons-material";
import { IconButton } from "@mui/material";
import { closeSnackbar, enqueueSnackbar } from "notistack";

export function SnackbarAction(snackbarId) {
    return (
        <IconButton onClick={() => closeSnackbar(snackbarId)} >
            <Close />
        </IconButton>
    );
}

export function noteSuccess(msg) {
    enqueueSnackbar(msg, { variant: "success" });
}

export function noteError(msg) {
    enqueueSnackbar(msg, { variant: "error" });
}

export function noteWarning(msg) {
    enqueueSnackbar(msg, { variant: "warning" });
}

export function noteInfo(msg) {
    enqueueSnackbar(msg, { variant: "info" });
}
