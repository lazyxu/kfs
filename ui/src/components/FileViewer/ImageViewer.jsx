import { Close, Download, Info } from "@mui/icons-material";
import { Box, Drawer, IconButton, Modal, Stack, Typography } from "@mui/material";
import { getSysConfig } from "hox/sysConfig";
import { useState } from "react";

export default function ({ open, setOpen, metadata, hash }) {
    const [info, setInfo] = useState(false);
    const sysConfig = getSysConfig().sysConfig;
    let { hash: hash2, exif, fileType } = metadata;
    return (
        <>
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
                    <Stack
                        sx={{
                            position: 'absolute',
                            right: 8,
                            top: 8,
                        }}
                        direction="row"
                        justifyContent="flex-end"
                        alignItems="flex-end"
                        spacing={0.5}
                    >
                        <IconButton
                            href={`${sysConfig.webServer}/api/v1/download?hash=${hash}`}
                            download="test.jpg"
                        >
                            <Download />
                        </IconButton>
                        <IconButton
                            onClick={() => setInfo(true)}
                        >
                            <Info />
                        </IconButton>
                        <IconButton
                            onClick={() => setOpen(false)}
                        >
                            <Close />
                        </IconButton>
                    </Stack>
                    <Box sx={{
                        width: "100%",
                        height: "100%",
                        display: "flex",
                        justifyContent: "center",
                        alignItems: "center",
                    }}
                    >
                        <img style={{ maxWidth: "100%", maxHeight: "100%" }} src={`${sysConfig.webServer}/api/v1/image?hash=${hash}`} />
                    </Box>
                </Box>
            </Modal>
            <Drawer
                anchor="right"
                open={info}
                onClose={() => setInfo(false)}
                sx={{ zIndex: 1350 }}
                SlideProps={{ sx: { maxWidth: "90%" } }}
            >
                <Box
                    sx={{ whiteSpace: "pre" }}
                    role="presentation"
                >
                    {JSON.stringify(metadata, null, 2)}
                </Box>
            </Drawer>
        </>
    );
}
