import { Box, Drawer } from "@mui/material";

export default function ({ hash, metadata, onClose }) {
    return (
        <Drawer
            anchor="right"
            open="true"
            onClose={onClose}
            sx={{ zIndex: 1350 }}
            SlideProps={{ sx: { maxWidth: "90%" } }}
        >
            <Box
                sx={{ whiteSpace: "pre" }}
            >
                {JSON.stringify(metadata, null, 2)}
            </Box>
        </Drawer>
    );
}
