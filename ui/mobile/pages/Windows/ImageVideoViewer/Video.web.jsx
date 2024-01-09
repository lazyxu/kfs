export default function ({ width, height, source }) {
    return (
        <video controls style={{ width, height }} data-setup='{}'>
            <source src={source} />
            您的浏览器不支持 HTML5 video 标签。
        </video>
    );
}
