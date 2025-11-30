export interface TestUser{
    username: string;
    password: string;
}

export const TEST_USERS: Record<string, TestUser> = {
    testuser: { username: 'testuser2', password: 'validpassword'},
    otherussr: { username: 'testuser3', password: '11111111'},
}
