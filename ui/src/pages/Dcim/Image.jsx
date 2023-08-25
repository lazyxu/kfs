import { Close } from "@mui/icons-material";
import { Box, IconButton, Modal, Typography } from "@mui/material";
import SvgIcon from "components/Icon/SvgIcon";
import { useState } from "react";
import styles from './image.module.scss';

export default function ({ metadata }) {
    const [open, setOpen] = useState(false);
    const handleOpen = () => setOpen(true);
    const handleClose = () => setOpen(false);
    let { hash, exif, fileType, shotTime, shotEquipment } = metadata;
    let time = shotTime.isValid() ? shotTime.format("YYYY年MM月DD日 HH:mm:ss") : "未知时间";
    let transform = "";
    switch (exif.Orientation) {
        case 2:
            transform = "rotateY(180deg)";
            break;
        case 3:
            transform = "rotate(180deg)";
            break;
        case 4:
            transform = "rotate(180deg)rotateY(180deg)";
            break;
        case 5:
            transform = "rotate(270deg)rotateY(180deg)";
            break;
        case 6:
            transform = "rotate(90deg)";
            break;
        case 7:
            transform = "rotate(90deg)rotateY(180deg)";
            break;
        case 8:
            transform = "rotate(270deg)";
            break;
    }
    return (
        <div>
            <img style={{ width: "100%", transform }} className={styles.clickable} src={`http://127.0.0.1:1123/thumbnail?size=256&cutSquare=true&hash=${hash}`} loading="lazy"
                title={time + "\n"
                    + shotEquipment + "\n"
                    + fileType.subType + "\n"
                    + (exif.GPSLatitudeRef ? (exif.GPSLatitudeRef + " " + exif.GPSLatitude + "\n") : "")
                    + (exif.GPSLongitudeRef ? (exif.GPSLongitudeRef + " " + exif.GPSLongitude + "\n") : "")
                    + hash}
                onClick={handleOpen}
            />
            <Modal
                open={open}
                onClose={handleClose}
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
                        onClick={handleClose}
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
                        <img style={{ maxWidth: "100%", maxHeight: "100%", transform }} src={`http://127.0.0.1:1123/api/v1/openDcim?hash=${hash}&fileType=${fileType.subType}`} />
                    </Box>
                </Box>
            </Modal>
        </div>
    );
}
