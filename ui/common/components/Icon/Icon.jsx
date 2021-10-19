import './index.module.scss';

export default function ({ icon, className }) {
  return (
    <svg
      aria-hidden="true"
      viewBox="0 0 200 200"
      preserveAspectRatio="xMinYMin meet"
      className={className}
    >
      <use xlinkHref={`#icon-${icon}`} />
    </svg>
  );
}
