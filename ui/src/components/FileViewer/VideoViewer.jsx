import { Close } from "@mui/icons-material";
import { Box, Button, IconButton, Modal } from "@mui/material";
import { getSysConfig } from "hox/sysConfig";
import { useRef } from "react";

export default function ({ open, setOpen, metadata, hash }) {
    const sysConfig = getSysConfig().sysConfig;
    let { hash: hash2, videoMetadata, fileType } = metadata;
    const playerRef = useRef(null);
    return (
        <Modal
            open={open}
            onClose={() => setOpen(false)}
            aria-labelledby="modal-modal-title"
            aria-describedby="modal-modal-description"
        >
            <Box sx={{
                width: "100%",
                height: "100%",
                backgroundColor: theme => theme.background.primary,
                color: theme => theme.context.primary
            }}
            >
                <IconButton
                    aria-label="close"
                    onClick={() => setOpen(false)}
                    sx={{
                        position: 'absolute',
                        right: 8,
                        top: 8,
                        color: (theme) => theme.palette.grey[500],
                    }}
                >
                    <Close />
                </IconButton>
                <Box sx={{
                    width: "100%",
                    height: "100%",
                    textAlign: "center"
                }}
                >
                    <video ref={playerRef} controls style={{ maxWidth: "100%", maxHeight: "100%" }} data-setup='{}'>
                        <source src={`${sysConfig.webServer}/api/v1/image?hash=${hash}`} />
                        您的浏览器不支持 HTML5 video 标签。
                    </video>
                    {/* <Button onClick={() => playerRef.current?.play()}>播放</Button>
                    <Button onClick={() => playerRef.current?.pause()}>暂停</Button> */}
                </Box>
            </Box>
        </Modal>
    );
}
