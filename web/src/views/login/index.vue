<template>
  <div class="content">
    <div class="form-wrapper">
      <div class="linker">
        <span class="ring"></span>
        <span class="ring"></span>
        <span class="ring"></span>
        <span class="ring"></span>
        <span class="ring"></span>
      </div>

      <form class="login-form">
        <input v-model="form.userName" type="text" name="username" :placeholder="$t('page.login.accout')" />
        <input v-model="form.password" type="password" name="password" :placeholder="$t('page.login.pass')" />
        <button type="submit" @click="onSubmit">{{ $t('page.login.submit') }}</button>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive } from 'vue';
import type { LoginInput } from '@/api/user/model';
import { $t } from '@/locales';
import useUserStore from '@/store/modules/user';
import { useRouter } from 'vue-router';
import { message } from '@/utils/message';

const userStore = useUserStore();
const router = useRouter();
const form = reactive<LoginInput>({
  userName: '',
  password: ''
});

const onSubmit = () => {
  userStore.Login(form).then((res) => {
    if (res.code === 200) {
      router.push({
        name: 'Task'
      });
    } else {
      message.error($t(''));
    }
  });
};
</script>

<style scoped>
.content {
  width: 680px;
  height: 420px;
  margin: 300px auto;
}

.form-wrapper {
  margin: 32px auto;
  width: 264px;
  height: 253px;
  position: relative;
  border: 1px solid rgb(197, 200, 204);
  background-color: rgb(248, 249, 250);
  text-align: center;
  border-radius: 5px;
  box-shadow: 0 1px 0 rgb(255, 255, 255), 0 2px 0 rgb(197, 200, 204), 0 3px 0 rgb(255, 255, 255),
    0 4px 0 rgb(197, 200, 204);
}

.form-wrapper::before {
  content: '';
  display: block;
  height: 37px;
  border-bottom: 1px solid rgb(197, 200, 204);
  border-radius: 5px 5px 0 0;
  box-shadow: inset 2px 2px 0 rgb(255, 255, 255);
}

.form-wrapper .login-form {
  padding-top: 40px;
  box-shadow: inset 2px 2px 0 rgb(255, 255, 255);
}

input[name='username'],
input[name='password'] {
  height: 40px;
  width: 200px;
  padding-left: 15px;
  display: block;
  border: 1px solid rgb(197, 200, 204);
  background-color: rgb(228, 230, 233);
  margin: 0 auto;
}

input[name='username'] {
  border-bottom: none;
  border-radius: 5px 5px 0 0;
  /* box-shadow: inset 0 0 5px 5px rgb(212, 214, 217); */
}

input[name='password'] {
  border-radius: 0 0 5px 5px;
}

button[type='submit'] {
  margin-top: 25px;
  width: 215px;
  height: 44px;
  color: #fff;
  font-size: 20px;
  border: none;
  background: linear-gradient(orange, red);
}

.linker {
  position: absolute;
  height: 40px;
  width: 248px;
  top: 18px;
  left: 10px;
}

/* 上环 */
.ring {
  position: relative;
  display: inline-block;
  border: 1px solid rgb(163, 164, 167);
  background-color: rgb(220, 222, 225);
  height: 12px;
  width: 12px;
  border-radius: 6px;
  margin-right: 33px;
}

.ring:last-child {
  margin-right: 0;
}

/* 下环 */
.ring::before {
  content: '';
  position: absolute;
  bottom: -25px;
  left: -1px;
  border: 1px solid rgb(163, 164, 167);
  background-color: rgb(220, 222, 225);
  height: 12px;
  width: 12px;
  border-radius: 6px;
}

/* 中间的竖线 */
.ring::after {
  content: '';
  position: absolute;
  top: 2px;
  left: 2px;
  width: 6px;
  height: 30px;
  border: 1px solid rgb(202, 202, 202);
  background-color: rgb(255, 255, 255);
  border-radius: 3px;
}
</style>
