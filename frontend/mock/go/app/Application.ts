import { resultSuccess, resultPageSuccess, resultError, getTimestamps } from '../../helper';
import * as Application from '/@wails/go/app/Application';
import { lowerCase, remove } from 'lodash-es';

const applications = [
  {
    id: 1,
    name: 'Example',
    describe: 'I am an example record',
    status: 0,
    ...getTimestamps(),
  },
  {
    id: 2,
    name: 'Test',
    describe: 'I am an test record',
    status: 1,
    ...getTimestamps(),
  },
];

let idCursor = applications.length;

export default {
  ApplicationList: (page: number, pageSize: number, name: string, status: number) => {
    let data = applications;
    if (name) data = applications.filter((x) => lowerCase(x.name) === lowerCase(name));
    if (status) data = applications.filter((x) => x.status === status);

    return resultPageSuccess(page, pageSize, data);
  },
  AddApplication: (application) => {
    idCursor++;

    const newApplication = {
      id: idCursor,
      status: 0,
      ...application,
      ...getTimestamps(),
    };

    applications.push(newApplication);

    return resultSuccess(true);
  },
  UpdateApplication: (application) => {
    const { id } = application;
    const foundApplication = applications.find((x) => x.id === id);

    if (!foundApplication) {
      return resultError('Record not found');
    }

    // Update application
    Object.assign(foundApplication, { ...application, ...getTimestamps('updatedAt') });

    return resultSuccess(true);
  },
  DeleteApplication: (id) => {
    remove(applications, (x) => x.id === id);
    return resultSuccess(true);
  },
  QueryApplicationById: (id) => {
    const foundApplication = applications.find((x) => x.id === id);
    return resultSuccess(foundApplication);
  },
} as Partial<typeof Application>;
