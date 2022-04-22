import { HelloWorldUserFactory } from './HelloWorldUser';
export enum name {
  HelloWorldUser = 'hello_world_users',
}

export const factory = {
  HelloWorldUserFactory,
};

export default { name, factory };
