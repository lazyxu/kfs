import { AllInbox, Close, Download, Info, PrivacyTip } from "@mui/icons-material";
import { Box, IconButton, Modal, Stack } from "@mui/material";
import { getSysConfig } from "hox/sysConfig";
import { useState } from "react";
import SameFiles from "./SameFiles";
import Metadata from "./Metadata";
import Attribute from "./Attribute";

export default function ({ open, setOpen, metadata, hash, attribute }) {
    const [openMetadata, setOpenMetadata] = useState(false);
    const [openSameFiles, setOpenSameFiles] = useState(false);
    const [openAttribute, setOpenAttribute] = useState(false);
    const [sameFiles, setSameFiles] = useState([]);
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
                            download={attribute ? attribute.name : sameFiles[0].name}
                        >
                            <Download />
                        </IconButton>
                        <IconButton
                            onClick={() => setOpenSameFiles(true)}
                        >
                            <AllInbox />
                        </IconButton>
                        {attribute && <IconButton
                            onClick={() => setOpenAttribute(true)}
                        >
                            <Info />
                        </IconButton>}
                        <IconButton
                            onClick={() => setOpenMetadata(true)}
                        >
                            <PrivacyTip />
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
            <SameFiles open={openSameFiles} setOpen={setOpenSameFiles} hash={hash} sameFiles={sameFiles} setSameFiles={setSameFiles} />
            <Metadata open={openMetadata} setOpen={setOpenMetadata} metadata={metadata} />
            {attribute && <Attribute open={openAttribute} setOpen={setOpenAttribute} attribute={attribute} />}
        </>
    );
}
