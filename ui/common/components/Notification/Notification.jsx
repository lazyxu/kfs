import { useSnackbar } from 'notistack';
import IconButton from "@mui/material/IconButton";
import CloseIcon from "@mui/material/SvgIcon/SvgIcon";
import Snackbar from "@mui/material/Snackbar/Snackbar";
import Slide from "@mui/material/Slide/Slide";
import React, { Fragment, useEffect, useState } from "react";

const useNotification = () => {
    const [conf, setConf] = useState({});
    const { enqueueSnackbar, closeSnackbar } = useSnackbar();
    // const content = (key, message) => (
    //     <Fragment>
    //         <Snackbar
    //             anchorOrigin={{ vertical: 'top', horizontal: 'right' }}
    //             // open={open}
    //             onClose={() => closeSnackbar(key)}
    //             message={message}
    //             key={key}
    //         />
    //     </Fragment>
    // );
    useEffect(() => {
        if (conf?.msg) {
            let variant = 'info';
            if (conf.variant) {
                variant = conf.variant;
            }
            enqueueSnackbar(conf.msg, {
                variant: variant,
                autoHideDuration: 5000,
                anchorOrigin: {
                    vertical: 'top',
                    horizontal: 'right',
                },
                TransitionComponent: Slide,
                // content,
            });
        }
    }, [conf]);
    const sendError = (e) => setConf({ msg: e.message, variant: 'error' })
    return [conf, setConf, sendError];
};

export default useNotification;