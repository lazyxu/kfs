import { Button, DialogActions, DialogContent, TextField } from "@mui/material";
import { useState, useEffect } from 'react';
import { createFolder } from "../../../api/fs";
import { noteSuccess, noteError } from "../../../components/Notification";

// 系统保留名称列表（不区分大小写）
const RESERVED_NAMES = [
  'CON', 'PRN', 'AUX', 'NUL', 
  'COM1', 'COM2', 'COM3', 'COM4', 'COM5', 'COM6', 'COM7', 'COM8', 'COM9',
  'LPT1', 'LPT2', 'LPT3', 'LPT4', 'LPT5', 'LPT6', 'LPT7', 'LPT8', 'LPT9'
];

export default function ({ onClose, dirPath }) {
  const [name, setName] = useState("新建文件夹");
  const [prem, setPrem] = useState("rwxrwxrwx");
  const [nameError, setNameError] = useState("");
  const [premError, setPremError] = useState("");
  
  // 名称验证逻辑
  const validateName = (value) => {
    if (!value) return "文件夹名称不能为空";
    if (value.length > 255) return "名称长度不能超过255个字符";
    
    // 检查非法字符[1,3](@ref)
    const forbiddenChars = /[\\/:*?"<>|]/;
    if (forbiddenChars.test(value)) {
      return "名称不能包含 \\ / : * ? \" < > | 等特殊字符";
    }
    
    // 检查保留名称[1](@ref)
    if (RESERVED_NAMES.includes(value.toUpperCase())) {
      return "不能使用系统保留名称";
    }
    
    return "";
  };

  // 权限验证逻辑（符合Linux权限格式）[9,10,11](@ref)
  const validatePrem = (value) => {
    if (value.length !== 9) return "权限必须是9个字符";
    
    const validChars = ['r', 'w', 'x', '-'];
    for (let i = 0; i < value.length; i++) {
      if (!validChars.includes(value[i])) {
        return `位置${i+1}的字符无效（只允许 r, w, x, -）`;
      }
    }
    
    // 验证每组权限格式（owner/group/others）
    const groups = [value.substring(0, 3), value.substring(3, 6), value.substring(6, 9)];
    for (const group of groups) {
      if (group.length !== 3) return "权限格式不正确";
    }
    
    return "";
  };

  // 实时验证输入
  useEffect(() => {
    setNameError(validateName(name));
    setPremError(validatePrem(prem));
  }, [name, prem]);

  const convertPremToDecimal = (prem) =>{
    // 验证权限字符串格式
    if (typeof prem !== 'string' || prem.length !== 9) {
        return -1; // 格式无效
    }

    // 将权限映射为二进制位
    let binary = '';
    for (let i = 0; i < 9; i++) {
        binary += (prem[i] === '-') ? '0' : '1'; // 有权限=1，无权限=0
    }

    // 将9位二进制转换为十进制整数
    return parseInt(binary, 2);
}

  // 提交处理
  const handleSubmit = async () => {
    const nameErr = validateName(name);
    const premErr = validatePrem(prem);
    
    if (nameErr || premErr) {
      setNameError(nameErr || "");
      setPremError(premErr || "");
      return;
    }
    
    try {
      await createFolder(dirPath, name, convertPremToDecimal(prem));
      noteSuccess("新建文件夹成功");
      onClose();
    } catch (err) {
      noteError("新建文件夹失败：" + (err?.message || err));
    }
  };

  return (
    <>
      <DialogContent>
        <TextField
          autoFocus
          margin="dense"
          placeholder="文件夹名字"
          fullWidth
          variant="outlined"
          value={name}
          error={!!nameError}
          helperText={nameError}
          onChange={(e) => setName(e.target.value)}
          FormHelperTextProps={{ style: { color: 'red' } }} // 红色错误提示[6,7](@ref)
        />
        <TextField
          placeholder="权限（格式: rwxrwxrwx）"
          margin="dense"
          fullWidth
          variant="outlined"
          value={prem}
          error={!!premError}
          helperText={premError || "格式：所有者|组|其他人的权限（r=读, w=写, x=执行）"}
          onChange={(e) => setPrem(e.target.value)}
          FormHelperTextProps={{ 
            style: { color: premError ? 'red' : 'inherit' } // 错误时红色提示[6,7](@ref)
          }}
        />
      </DialogContent>
      <DialogActions>
        <Button 
          variant="outlined" 
          sx={{ width: "10em" }} 
          disabled={!!nameError || !!premError} // 有错误时禁用按钮[2](@ref)
          onClick={handleSubmit}
        >
          确定
        </Button>
      </DialogActions>
    </>
  );
}