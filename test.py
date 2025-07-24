import os
import shutil
import argparse

def transfer_files(src_root, dst_root):
    """
    将文件从旧目录结构迁移到新目录结构
    :param src_root: 源根目录 (tmp-07241512)
    :param dst_root: 目标根目录 (tmp)
    """
    # 遍历源目录中的所有文件
    for size_dir in os.listdir(src_root):
        src_size_path = os.path.join(src_root, size_dir)
        
        # 确保是目录且包含数字（文件大小）
        if not os.path.isdir(src_size_path) or not size_dir.isdigit():
            continue
            
        size = int(size_dir)
        mod_dir = str(size % 1024)  # 计算取模目录
        
        entry_count = len(os.listdir(src_size_path))
        if entry_count > 1:
            print(f"目录: {src_size_path}, 条目数量: {entry_count}")
        # 处理每个哈希文件
        for file_name in os.listdir(src_size_path):
            src_file_path = os.path.join(src_size_path, file_name)
            
            # 跳过子目录
            if not os.path.isfile(src_file_path):
                continue
                
            # 解析哈希值
            if '_' not in file_name:
                print(f"跳过无效文件名: {file_name}")
                continue
                
            md5_hash, sha256_hash = file_name.split('_', 1)
            
            # 构建目标路径
            dst_dir = os.path.join(dst_root, mod_dir)
            dst_file_name = f"{size_dir}_{md5_hash}_{sha256_hash}"
            dst_file_path = os.path.join(dst_dir, dst_file_name)
            
            # 创建目标目录（如果不存在）
            os.makedirs(dst_dir, exist_ok=True)
            
            # 移动文件
            try:
                shutil.move(src_file_path, dst_file_path)
                print(f"移动成功: {src_file_path} → {dst_file_path}")
            except Exception as e:
                print(f"移动失败 {src_file_path}: {str(e)}")

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description='迁移文件目录结构')
    parser.add_argument('--src', default='tmp-07241512', help='源根目录路径')
    parser.add_argument('--dst', default='tmp', help='目标根目录路径')
    
    args = parser.parse_args()
    
    print(f"开始迁移: {args.src} → {args.dst}")
    transfer_files(args.src, args.dst)
    print("迁移完成")