<template>
  <div
    class="flex bg-yellow-500 dark:bg-yellow-900 justify-center flex-col min-h-screen min-w-full rounded-lg"
  >
    <div class="bg-white dark:bg-black shadow-md mx-auto rounded-lg p-6 m-5">
      <h1 class="text-3xl leading-9 font-extrabold mb-8">{{ msg }}</h1>
      <input
        type="email"
        class="md:flex p-2 mb-3 w-full text-xl border border-yellow-400 rounded"
        placeholder="Email Id"
        v-model="this.email"
      />
     

      <h2 class="self-start text-left mb-3 text-xl">Topics</h2>

      <span v-for="topic of this.topics" class="mx-3" :key="topic">
        <input
          class="row-gap-4 "
          type="checkbox"
          :value="topic"
          :id="topic"
          v-model="this.checkedTopics"
        />
        <label :for="topic">
          <span class="mx-2 text-lg">{{ topic }}</span>
        </label>
      </span>

      <div>
        <button
          class="bg-yellow-400 hover:shadow-md hover:bg-yellow-600 hover:-translate-y-1 ease-in-out duration-75 transform hover:scale-125 text-white font-bold py-2 px-4 m-5 rounded text-lg"
          @click="submit()"
        >
          Submit
        </button>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "HelloWorld",
  props: {
    msg: String,
  },
  data() {
    return {
      email: "",
      topics: ["Front", "Jobs", "Ask", "Show"],
      checkedTopics: [],
    };
  },
  methods: {
    submit() {
      console.log(this.email);
      console.log(Object.values(this.checkedTopics));
    },
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
