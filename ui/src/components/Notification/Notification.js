import { enqueueSnackbar } from "notistack"

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
