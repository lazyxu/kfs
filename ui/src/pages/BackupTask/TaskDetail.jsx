import { Close } from "@mui/icons-material";
import { Box, IconButton, Typography } from "@mui/material";

export default function ({ taskDetail, setTaskDetail }) {
    return (
        <Box sx={{
            width: "100%", flex: "1",
            display: 'flex', flexDirection: 'column', minHeight: '0'
        }}>
            <Typography variant="h6" noWrap component="div" sx={{ marginLeft: 2 }}>
                任务日志： {taskDetail}
                <IconButton
                    aria-label="close"
                    onClick={() => setTaskDetail()}
                    sx={{
                        color: (theme) => theme.palette.grey[500],
                    }}
                >
                    <Close />
                </IconButton>
            </Typography>
        </Box>
    );
}
