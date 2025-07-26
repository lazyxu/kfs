import useResourceManager, { openDir, openDrivers } from '@kfs/common/hox/resourceManager';
import SvgIcon from "@kfs/mui/components/Icon/SvgIcon";
import { ArrowBack, ArrowForward, ArrowForwardIos, ArrowUpward, Refresh, Search } from "@mui/icons-material";
import { Button, IconButton, Stack, TextField } from "@mui/material";
import { useTheme } from "@mui/material/styles";
import { useRef, useState } from "react";

function PathButton({ name, icon, dirPath }) {
    const [resourceManager, setResourceManager] = useResourceManager();
    return (
        <>
            {icon && <Button sx={{minWidth: "28px"}}
                onClick={e => {
                    e.preventDefault(); e.stopPropagation();
                    openDir(setResourceManager, dirPath);
                }}
                >
                <SvgIcon icon={icon} style={{ height: "24px", width: "24px" }} />
            </Button>}
            {name && <Button
                variant="text"
                size="small"
                sx={{
                    height: "32px",
                    minWidth: 0,
                    display: "flex",
                    alignItems: "center",
                }}
                onClick={e => {
                    e.preventDefault(); e.stopPropagation();
                    openDir(setResourceManager, dirPath);
                }}
            >
                <span style={{ lineHeight: "normal" }}>{name}</span>
            </Button>}
            <IconButton
                size="small"
                sx={{
                    minWidth: 0,
                    boxShadow: "none",
                    borderLeft: "none",
                    m: 0,
                    width: 24,
                    height: 24,
                    display: "flex",
                    alignItems: "center",
                    justifyContent: "center",
                }}
            >
                <ArrowForwardIos fontSize="small" />
            </IconButton>
        </>
    );
}

export default () => {
    const [resourceManager, setResourceManager] = useResourceManager();
    const dirPath = resourceManager.dirPath || [];
    const [editMode, setEditMode] = useState(false);
    const [inputValue, setInputValue] = useState("");
    const inputRef = useRef();
    const theme = useTheme();
    const [searchMode, setSearchMode] = useState(false);
    const [searchValue, setSearchValue] = useState("");
    // 历史记录
    const [history, setHistory] = useState([dirPath]);
    const [historyIndex, setHistoryIndex] = useState(0);

    // 跳转目录并维护历史
    const jumpDir = (newPath) => {
        openDir(setResourceManager, newPath);
        setHistory(prev => {
            const updated = prev.slice(0, historyIndex + 1).concat([newPath]);
            setHistoryIndex(updated.length - 1);
            return updated;
        });
    };

    // 后退
    const handleBack = () => {
        if (historyIndex > 0) {
            setHistoryIndex(historyIndex - 1);
            openDir(setResourceManager, history[historyIndex - 1]);
        }
    };

    // 前进
    const handleForward = () => {
        if (historyIndex < history.length - 1) {
            setHistoryIndex(historyIndex + 1);
            openDir(setResourceManager, history[historyIndex + 1]);
        }
    };

    // 上级
    const handleUp = () => {
        if (dirPath.length > 0) {
            const upPath = dirPath.slice(0, -1);
            jumpDir(upPath);
        }
    };

    // 刷新
    const handleRefresh = () => {
        openDir(setResourceManager, dirPath, { refresh: true });
    };

    // 搜索
    const handleSearch = () => {
        // TODO: 实现搜索逻辑
        setSearchMode(false);
        // openDir/setResourceManager等，带搜索参数
    };

    // 编辑路径
    const handleEdit = () => {
        setInputValue(dirPath.join("\\"));
        setEditMode(true);
        setTimeout(() => inputRef.current?.focus(), 100);
    };
    const handleInputChange = (e) => setInputValue(e.target.value);
    const handleInputBlur = () => setEditMode(false);
    const handleInputKeyDown = (e) => {
        if (e.key === "Enter") {
            const parts = inputValue.split("/");
            console.log("Navigating to:", parts);
            if (parts[0]?.length > 0) {
                openDir(setResourceManager, parts.slice(1));
            } else {
                openDrivers(setResourceManager);
            }
            setEditMode(false);
        }
    };

    // Win11风格分隔符
    const Separator = () => (
        <ArrowForwardIos />
    );

    return (
        <Stack
            sx={{
                padding: "0 1em",
                margin: "0.3em",
            }}
            direction="row"
            alignItems="center"
            spacing={0.5}
        >
            {/* Win11风格左侧操作按钮区 */}
            <Stack direction="row" alignItems="center" spacing={0.5} >
                <IconButton size="small" onClick={handleBack} disabled={historyIndex === 0} title="后退">
                    <ArrowBack />
                </IconButton>
                <IconButton size="small" onClick={handleForward} disabled={historyIndex === history.length - 1} title="前进">
                    <ArrowForward />
                </IconButton>
                <IconButton size="small" onClick={handleUp} disabled={dirPath.length === 0} title="上级">
                    <ArrowUpward />
                </IconButton>
                <IconButton size="small" onClick={handleRefresh} title="刷新">
                    <Refresh />
                </IconButton>
            </Stack>
            {/* 中间地址栏 */}
            <Stack
                direction="row"
                alignItems="center"
                spacing={0.5}
                sx={{ flex: 1, minWidth: 0, position: "relative" }}
                onClick={e => {
                    if (!editMode) {
                        setInputValue((resourceManager.driver?.name || "") + "/" + dirPath.join("/"));
                        setEditMode(true);
                        setTimeout(() => inputRef.current?.focus(), 100);
                    }
                }}
            >
                {editMode ? (
                    <TextField
                        inputRef={inputRef}
                        value={inputValue}
                        onChange={handleInputChange}
                        onBlur={handleInputBlur}
                        onKeyDown={handleInputKeyDown}
                        size="small"
                        sx={{
                            width: "100%",
                            bgcolor: theme.palette.action.hover,
                            borderRadius: 1,
                            flex: 1,
                        }}
                        autoFocus
                    />
                ) : (
                    <Stack
                        direction="row"
                        alignItems="center"
                        spacing={0.5}
                        sx={{
                            flex: 1,
                            height: '38px', // TextField高度
                            borderRadius: 1,
                            paddingLeft: "0.5em",
                            border: "1px solid var(--mui-palette-Chip-defaultBorder)",
                            "&:hover": { cursor: "pointer", borderColor: "var(--mui-palette-action-active)" },
                        }}
                    >
                        <PathButton icon="wangpan" dirPath={[]} />
                        {dirPath.map((name, i) => (
                            <Stack direction="row" alignItems="center" key={i}>
                                <PathButton
                                    name={name}
                                    dirPath={dirPath.slice(0, i + 1)}
                                />
                            </Stack>
                        ))}
                    </Stack>
                )}
            </Stack>
            {/* 右侧搜索栏，始终显示搜索按钮，点击后显示输入框 */}
            <Stack direction="row" alignItems="center" spacing={0.5} sx={{ minWidth: 200, flex: "0 0 320px" }}>
                <TextField
                    value={searchValue}
                    onChange={e => setSearchValue(e.target.value)}
                    onKeyDown={e => {
                        if (e.key === "Enter") handleSearch();
                    }}
                    size="small"
                    sx={{
                        minWidth: 160,
                        maxWidth: 320,
                        bgcolor: theme.palette.action.hover,
                        borderRadius: 1,
                        flex: 1,
                        mr: 1,
                    }}
                    placeholder="搜索"
                    InputProps={{
                        endAdornment: (
                            <IconButton
                                size="small"
                                onClick={handleSearch}
                                title="搜索"
                                sx={{
                                    p: 0.5,
                                    bgcolor: theme.palette.action.hover,
                                    "&:hover": { bgcolor: theme.palette.action.selected }
                                }}
                            >
                                <Search />
                            </IconButton>
                        )
                    }}
                />
            </Stack>
        </Stack>
    );
};
