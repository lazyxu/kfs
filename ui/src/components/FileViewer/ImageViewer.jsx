import { Close } from "@mui/icons-material";
import { Box, IconButton, Modal } from "@mui/material";
import { getTransform } from "api/utils/api";

export default function ({ open, setOpen, metadata, hash }) {
    let { hash: hash2, exif, fileType } = metadata;
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
                    <img style={{ maxWidth: "100%", maxHeight: "100%", transform: getTransform(exif.Orientation) }} src={`${location.origin}/api/v1/image?hash=${hash}`} />
                </Box>
            </Box>
        </Modal>
    );
}
