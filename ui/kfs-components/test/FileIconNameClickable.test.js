import { mount } from 'enzyme';

import FileIconNameClickable from '../components/FileIconNameClickable';

describe('Enzyme Mount', function () {
  it('Delete Todo', function () {
    let app = mount(<FileIconNameClickable />);
    let todoLength = app.find('li').length;
    app.find('button.delete').at(0).simulate('click');
    expect(app.find('li').length).to.equal(todoLength - 1);
  });
});
