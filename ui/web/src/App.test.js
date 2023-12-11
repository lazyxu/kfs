import { render, screen } from '@testing-library/react';
import { HoxRoot } from "hox";
import App from './App';

test('renders src/App.js', () => {
  render(
    <HoxRoot>
      <App />
    </HoxRoot>);
  const element = screen.getByText(/文件/i);
  expect(element).toBeInTheDocument();
});
