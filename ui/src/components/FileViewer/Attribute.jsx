import { Box, Drawer, } from "@mui/material";

export default function ({ attribute, open, setOpen }) {
    return (
        <Drawer
            anchor="right"
            open={open}
            onClose={() => setOpen(false)}
            sx={{ zIndex: 1350 }}
            SlideProps={{ sx: { maxWidth: "90%" } }}
        >
            <Box
                sx={{ whiteSpace: "pre" }}
            >
                {JSON.stringify(attribute, null, 2)}
            </Box>
        </Drawer>
    );
}
