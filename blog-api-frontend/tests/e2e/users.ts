export interface TestUser{
    username: string;
    password: string;
}

export const TEST_USERS: Record<string, TestUser> = {
    testuser: { username: 'e2e_test', password: 'e2e_password'},
    otheruser: { username: 'e2e_test2', password: 'e2e_password2'},
}
